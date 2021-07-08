FROM public.ecr.aws/amazonlinux/amazonlinux:2

RUN yum update -y \
    && yum install -y \
    # for ssm agent download
    wget \
    # for NonInteractiveCommand
    curl \
    # for useradd
    # TODO: we probably don't need this as container mode doesn't create ssm-user.
    shadow-utils \
    # Dependency of session manager logging
    # TODO: remove it since we don't really want logging.
    util-linux \
    && wget https://s3.eu-north-1.amazonaws.com/amazon-ssm-eu-north-1/latest/linux_amd64/amazon-ssm-agent.rpm \
    && yum remove -y \
    wget \
    && yum clean all \
    && rpm -i ./amazon-ssm-agent.rpm \
    && rm -f ./amazon-ssm-agent.rpm

RUN echo "{\"Agent\": {\"ContainerMode\": true},\"Identity\": {\"ConsumptionOrder\": [\"OnPrem\"]}}" > /etc/amazon/ssm/amazon-ssm-agent.json

# Debug log config
RUN echo "<!--amazon-ssm-agent uses seelog logging --> \
<!--Seelog has github wiki pages, which contain detailed how-tos references: https://github.com/cihub/seelog/wiki --> \
<!--Seelog examples can be found here: https://github.com/cihub/seelog-examples --> \
<seelog type=\"adaptive\" mininterval=\"2000000\" maxinterval=\"100000000\" critmsgcount=\"500\" minlevel=\"trace\"> \
    <exceptions> \
        <exception filepattern=\"test*\" minlevel=\"error\"/> \
    </exceptions> \
    <outputs formatid=\"fmtinfo\"> \
        <console formatid=\"fmtinfo\"/> \
        <rollingfile type=\"size\" filename=\"/var/log/amazon/ssm/amazon-ssm-agent.log\" maxsize=\"30000000\" maxrolls=\"5\"/> \
        <filter levels=\"error,critical\" formatid=\"fmterror\"> \
            <rollingfile type=\"size\" filename=\"/var/log/amazon/ssm/errors.log\" maxsize=\"10000000\" maxrolls=\"5\"/> \
        </filter> \
    </outputs> \
    <formats> \
        <format id=\"fmterror\" format=\"%Date %Time %LEVEL [%FuncShort @ %File.%Line] %Msg%n\"/> \
        <format id=\"fmtdebug\" format=\"%Date %Time %LEVEL [%FuncShort @ %File.%Line] %Msg%n\"/> \
        <format id=\"fmtinfo\" format=\"%Date %Time %LEVEL %Msg%n\"/> \
    </formats> \
</seelog>" > /etc/amazon/ssm/seelog.xml
 
CMD ["/usr/bin/amazon-ssm-agent"]
import React, { Component, PropTypes } from 'react';

import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import Button from 'components/buttons/Button';
import helpers from 'components/forms/RegistrationForm/KolideDetails/helpers';
import InputFieldWithIcon from 'components/forms/fields/InputFieldWithIcon';

const formFields = ['kolide_web_address'];
const { validate } = helpers;

class KolideDetails extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      kolide_web_address: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func.isRequired,
  };

  render () {
    const { fields, handleSubmit } = this.props;

    return (
      <div>
        <InputFieldWithIcon
          {...fields.kolide_web_address}
          placeholder="Kolide Web Address"
        />
        <Button
          onClick={handleSubmit}
          text="Submit"
          variant="gradient"
        />
      </div>
    );
  }
}

export default Form(KolideDetails, {
  fields: formFields,
  validate,
});